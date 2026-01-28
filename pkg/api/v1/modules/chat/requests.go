package chat

type ListChatMessagesRequest struct {
	Limit  int32 `json:"limit" form:"limit,default=10" example:"10" query:"limit"`
	Offset int32 `json:"offset" form:"offset,default=0" example:"0" query:"offset"`
}

type SendChatMessageRequest struct {
	Content              string `json:"content" binding:"required" example:"Hello, how are you?"`
	Code                 string `json:"code" binding:"required" example:"print('Hello, world!')"`
	ExerciseInstructions string `json:"exercise_instructions" binding:"required" example:"Write a function that prints 'Hello, world!'"`
}
