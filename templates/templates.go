package templates

import "html/template"

var General = template.Must(template.ParseGlob("templates/*.html"))
var FormTemplate = template.Must(template.ParseFiles("templates/form/edit_task_form.html"))
