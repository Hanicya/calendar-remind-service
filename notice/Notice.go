package notice

type Notice interface {
	// Notice 通知方法 contactInfo 通知方式 context 通知内容
	Notice(contactInfo string, context string)
}
