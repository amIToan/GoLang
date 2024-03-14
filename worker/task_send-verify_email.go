package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/util"
	"sgithub.com/techschool/simplebank/validator"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	paresedJsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, paresedJsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", taskInfo.Queue).Int("max-retry", taskInfo.MaxRetry).Msg("Task enqueued")
	return nil
}
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", asynq.SkipRetry)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("the user doesn't exist : %w", err)
		}
		return fmt.Errorf("failed to get the user : %w", err)
	}
	if err = validator.ValidateEmail(user.Email); err != nil {
		return fmt.Errorf("failed to insert verifying Email : %w", err)
	}
	verifiedEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomStr(32),
	})
	if err != nil {
		return fmt.Errorf("failed to insert verifying Email : %w", err)
	}
	subject := "A test email"
	toVerifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?id=%d&secret_code=%s", verifiedEmail.ID, verifiedEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, toVerifyUrl)
	to := []string{user.Email}
	processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed Task")

	return nil
}
