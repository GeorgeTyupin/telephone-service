package mapper

import (
	"fmt"
	"telephone-service/codegen/oas"
	"telephone-service/internal/model"
)

func CallFromRequest(
	call api.CreateCallRequest,
	first api.Subscriber,
	second api.Subscriber,
) (model.Call, error) {
	firstSubscriber, err := SubscriberFromRequest(first)
	if err != nil {
		return model.Call{}, fmt.Errorf("преобразование звонка, first_abonent: %w", err)
	}

	secondSubscriber, err := SubscriberFromRequest(second)
	if err != nil {
		return model.Call{}, fmt.Errorf("преобразование звонка, second_abonent: %w", err)
	}

	duration, err := model.StringToDuration(call.Duration)
	if err != nil {
		return model.Call{}, fmt.Errorf("преобразование звонка, duration: %w", err)
	}

	createdAt, err := model.StringToTime(call.CreatedAt)
	if err != nil {
		return model.Call{}, fmt.Errorf("преобразование звонка, created_at: %w", err)
	}

	return model.Call{
		Id:               call.ID,
		FirstSubscriber:  firstSubscriber,
		SecondSubscriber: secondSubscriber,
		Duration:         duration,
		CreatedAt:        createdAt,
	}, nil
}
