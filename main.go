package main

import (
	"fmt"
	"note_api/http"
	"note_api/note"
)

func main()  {
	repo := note.NewRepository()
	httpHandler := http.NewHTTPHandlers(repo)
	httpServer := http.NewHTTPServer(httpHandler)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server:", err)
	}
}

// todo: fix когда обновляю сразу название и содержимое заметки, отправляю запрос, пишется ошибка 404, при этом когда вывожу заметку по старому названию тоже 404, когда по новому, то она находится, но без изменения контента. то есть когда меняю все - меняется ток название , при этом выводя ошибку а не измененную заметку