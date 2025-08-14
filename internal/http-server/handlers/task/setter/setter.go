package setter

import (
	"log/slog"
	"net/http"
	respmodel "test-task-lo/internal/http-server/model/response"
	logext "test-task-lo/internal/lib/log"

	"github.com/go-chi/render"
)

type Request struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Response struct {
	respmodel.Response
	ID string `json:"id"`
}

type TaskSetter interface {
	SetTask(name string, desc string) (string, error)
}

func NewSet(log logext.Logger, taskSetter TaskSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.set.new"

		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("can't decode request body", slog.String("op", op), logext.Err(err))
			render.JSON(w, r, respmodel.NewStatusError("can't decode request body"))
			return
		}

		//TODO: use the validator lib for validation
		if req.Name == "" {
			log.Error("invalid reqest body. 'name' required", slog.String("op", op))
			render.JSON(w, r, respmodel.NewStatusError("cinvalid reqest body. 'name' required"))
			return
		}

		id, err := taskSetter.SetTask(req.Name, req.Desc)
		if err != nil {
			log.Error("failed to add task", slog.String("op", op), logext.Err(err))
			render.JSON(w, r, respmodel.NewStatusError("failed to add url"))
			return
		}
		log.Info("id is added", slog.String("op", op), slog.String("id", id))
		ResponseOK(w, r, id)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Response: respmodel.NewStatusOK(),
		ID:       id,
	})
}
