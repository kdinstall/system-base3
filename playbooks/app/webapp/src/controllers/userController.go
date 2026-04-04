package controllers

import (
	"net/http"
	"strconv"
	"web-sqlx-sqlite-user/src/db"
	tmpl "web-sqlx-sqlite-user/src/lib/template"

	"github.com/gin-gonic/gin"
)

// UserController はユーザ管理の各アクションを提供する
type UserController struct{}

// Index はユーザ一覧を表示する (GET /users)
func (uc *UserController) Index(c *gin.Context) {
	userDb := db.NewUserDb()
	users := userDb.ListUsers()

	c.HTML(http.StatusOK, "users.html", tmpl.MergeData(gin.H{
		"page_title": "ユーザ一覧",
		"users":      users,
		"flash":      c.Query("flash"),
	}))
}

// New は新規ユーザ作成フォームを表示する (GET /users/new)
func (uc *UserController) New(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", tmpl.MergeData(gin.H{
		"page_title": "ユーザ新規作成",
		"errors":     nil,
		"input":      gin.H{},
	}))
}

// Create は新規ユーザを登録する (POST /users)
func (uc *UserController) Create(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	// バリデーション
	errs := validateUserInput(name, email)
	if len(errs) > 0 {
		c.HTML(http.StatusUnprocessableEntity, "user_new.html", tmpl.MergeData(gin.H{
			"page_title": "ユーザ新規作成",
			"errors":     errs,
			"input":      gin.H{"name": name, "email": email},
		}))
		return
	}

	userDb := db.NewUserDb()
	if err := userDb.CreateUser(name, email); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "user_new.html", tmpl.MergeData(gin.H{
			"page_title": "ユーザ新規作成",
			"errors":     []string{"登録に失敗しました: " + err.Error()},
			"input":      gin.H{"name": name, "email": email},
		}))
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?flash=created")
}

// Edit はユーザ編集フォームを表示する (GET /users/:id/edit)
func (uc *UserController) Edit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", tmpl.MergeData(gin.H{"page_title": "Not Found"}))
		return
	}

	userDb := db.NewUserDb()
	user := userDb.GetUser(id)
	if user == nil {
		c.HTML(http.StatusNotFound, "404.html", tmpl.MergeData(gin.H{"page_title": "Not Found"}))
		return
	}

	c.HTML(http.StatusOK, "user_edit.html", tmpl.MergeData(gin.H{
		"page_title": "ユーザ編集",
		"user":       user,
		"errors":     nil,
	}))
}

// Update はユーザ情報を更新する (POST /users/:id)
func (uc *UserController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", tmpl.MergeData(gin.H{"page_title": "Not Found"}))
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")

	// バリデーション
	errs := validateUserInput(name, email)
	if len(errs) > 0 {
		userDb := db.NewUserDb()
		user := userDb.GetUser(id)
		c.HTML(http.StatusUnprocessableEntity, "user_edit.html", tmpl.MergeData(gin.H{
			"page_title": "ユーザ編集",
			"user":       user,
			"errors":     errs,
		}))
		return
	}

	userDb := db.NewUserDb()
	if err := userDb.UpdateUser(id, name, email); err != nil {
		user := userDb.GetUser(id)
		c.HTML(http.StatusUnprocessableEntity, "user_edit.html", tmpl.MergeData(gin.H{
			"page_title": "ユーザ編集",
			"user":       user,
			"errors":     []string{"更新に失敗しました: " + err.Error()},
		}))
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?flash=updated")
}

// Delete はユーザを削除する (POST /users/:id/delete)
func (uc *UserController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", tmpl.MergeData(gin.H{"page_title": "Not Found"}))
		return
	}

	userDb := db.NewUserDb()
	if err := userDb.DeleteUser(id); err != nil {
		c.HTML(http.StatusInternalServerError, "users.html", tmpl.MergeData(gin.H{
			"page_title": "ユーザ一覧",
			"users":      userDb.ListUsers(),
			"error":      "削除に失敗しました: " + err.Error(),
		}))
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?flash=deleted")
}

// validateUserInput はユーザ入力のバリデーションを行う
func validateUserInput(name, email string) []string {
	var errs []string
	if name == "" {
		errs = append(errs, "名前は必須です。")
	}
	if len([]rune(name)) > 100 {
		errs = append(errs, "名前は100文字以内で入力してください。")
	}
	if email == "" {
		errs = append(errs, "メールアドレスは必須です。")
	}
	if len(email) > 255 {
		errs = append(errs, "メールアドレスは255文字以内で入力してください。")
	}
	return errs
}
