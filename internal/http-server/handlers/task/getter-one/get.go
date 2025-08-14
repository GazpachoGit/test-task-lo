package getterone

import (
	"errors"
	"log/slog"
	"net/http"
	model "test-task-lo/internal/domain/models"
	respmodel "test-task-lo/internal/http-server/model/response"
	logext "test-task-lo/internal/lib/log"
	"test-task-lo/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Response struct {
	respmodel.Response
	Task model.Task
}

type TaskGetter interface {
	GetTask(id string) (model.Task, error)
}

func New(log logext.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get-one.new"

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("task id param is empty", slog.String("op", op))
			render.JSON(w, r, respmodel.NewStatusError("invalid request"))
			return
		}

		task, err := taskGetter.GetTask(id)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				log.Info("task not found", slog.String("op", op), slog.String("id", id))
				render.JSON(w, r, respmodel.NewStatusError("task not found"))
			} else {
				log.Error("failed to get task", slog.String("op", op), logext.Err(err))
				render.JSON(w, r, respmodel.NewStatusError("failed to get task"))
			}
			return
		}

		log.Info("url found", slog.String("op", op), slog.String("task", task.Name))
		ResponseOK(w, r, task)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, task model.Task) {
	render.JSON(w, r, Response{
		Response: respmodel.NewStatusOK(),
		Task:     task,
	})
}
