package gettermulty

import (
	"errors"
	"log/slog"
	"net/http"
	model "test-task-lo/internal/domain/models"
	respmodel "test-task-lo/internal/http-server/model/response"
	logext "test-task-lo/internal/lib/log"
	"test-task-lo/internal/storage"

	"github.com/go-chi/render"
)

type Response struct {
	respmodel.Response
	Tasks []model.Task `json:"tasks"`
}

type TaskGetter interface {
	GetTasks(status string) ([]model.Task, error)
}

func New(log logext.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get-all.new"

		status := r.URL.Query().Get("status")
		tasks, err := taskGetter.GetTasks(status)
		if err != nil {
			if errors.Is(err, storage.ErrTaskStatusInvalid) {
				log.Info("invalid task status", slog.String("op", op), slog.String("status", status))
				render.JSON(w, r, respmodel.NewStatusError("invalid task status"))
			} else {
				log.Error("failed to get tasks", slog.String("op", op), logext.Err(err))
				render.JSON(w, r, respmodel.NewStatusError("failed to get tasks"))
			}
			return
		}
		if len(tasks) == 0 {
			log.Info("no tasks found", slog.String("op", op))
		} else {
			log.Info("tasks found", slog.String("op", op), slog.Int("count", len(tasks)))
		}
		ResponseOK(w, r, tasks)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, tasks []model.Task) {
	render.JSON(w, r, Response{
		Response: respmodel.NewStatusOK(),
		Tasks:    tasks,
	})
}
