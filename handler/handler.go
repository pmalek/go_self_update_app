package handler

import (
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/pkg/errors"
	"github.com/pmalek/proton_task/update"
)

const (
	pageTemplate = `
<!DOCTYPE html>
<html>
<head>
<title> Server {{.Version}} </title> </head>
<body>

<h1>This server is version {{.Version}}</h1>
<a href="check">Check for new version</a>

<br>

{{if .NewVersion}}
New version is available: {{.NewVersion}} | <a href="install">Upgrade</a>
{{end}}

{{if .Error}}
Error updating server: {{.Error}}
{{end}}

</body>
</html>
`

	pageUpdateTemplate = `
<!DOCTYPE html>
<html>
<head>
<title> Updating... </title>
</head>

<body>
<h1> Updating server from version {{.Version}} to version {{.NewVersion}} ... </h1>
</body>

<script>
setTimeout(function(){ window.location.replace("/"); }, 2000);
</script> 
</html>
`
)

type status struct {
	Version    int
	NewVersion int
	Error      error
}

type Handler struct {
	template       *template.Template
	updateTemplate *template.Template
	muStatus       sync.RWMutex
	status         status
	updateProvider update.Provider
	ch             chan status
}

func New(version int, updateProvider update.Provider) (*Handler, error) {
	template, err := template.New("page").Parse(pageTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page template")
	}

	updateTemplate, err := template.New("update").Parse(pageUpdateTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page template")
	}

	h := &Handler{
		template:       template,
		updateTemplate: updateTemplate,
		status: status{
			Version:    version,
			NewVersion: 0,
		},
		updateProvider: updateProvider,
		ch:             make(chan status),
	}

	go h.watchUpdates()

	return h, nil
}

func (h *Handler) watchUpdates() {
	for s := range h.ch {
		if s.Error != nil {
			h.status.Error = s.Error
			continue
		}

		if s.NewVersion > h.status.Version {
			h.muStatus.Lock()
			h.status.NewVersion = s.NewVersion
			h.muStatus.Unlock()
			continue
		}
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	st := h.getStatus()
	if err := h.template.Execute(w, st); err != nil {
		log.Printf("ERROR: failed to render page: %v", err)
	}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	st := h.getStatus()

	newv, err := h.updateProvider.IsUpdateAvailable(st.Version)
	if err != nil {
		http.Error(w, "failed to check for update", http.StatusInternalServerError)
		h.ch <- status{
			Error: err,
		}
		return
	}

	h.ch <- status{
		NewVersion: newv,
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) Install(w http.ResponseWriter, r *http.Request) {
	st := h.getStatus()
	v := st.NewVersion

	if err := h.updateTemplate.Execute(w, h.status); err != nil {
		log.Printf("ERROR: failed to render update page: %v", err)
	}

	go func(ch chan status) {
		if err := h.updateProvider.Update(v); err != nil {
			ch <- status{
				Error: err,
			}
		}
	}(h.ch)
}

func (h *Handler) getStatus() status {
	h.muStatus.RLock()
	s := h.status
	h.muStatus.RUnlock()
	return s
}
