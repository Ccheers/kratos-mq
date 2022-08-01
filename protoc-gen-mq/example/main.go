package main

import (
	"context"
	"net/http"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/protoc-gen-mq/example/api/product/app/ecode"
	v1 "github.com/Ccheers/kratos-mq/protoc-gen-mq/example/api/product/app/v1"
)

type service struct{}

func (s service) MQ_CreateArticle(ctx context.Context, article *v1.Article) error {
	if article.AuthorId < 1 {
		return ecode.Errorf(http.StatusBadRequest, 400, "author id must > 0")
	}
	return nil
}

func (s service) CreateArticle(ctx context.Context, article *v1.Article) (*v1.Article, error) {
	if article.AuthorId < 1 {
		return nil, ecode.Errorf(http.StatusBadRequest, 400, "author id must > 0")
	}
	return article, nil
}

func (s service) GetArticles(ctx context.Context, req *v1.GetArticlesReq) (*v1.GetArticlesResp, error) {
	if req.AuthorId < 0 {
		return nil, ecode.Errorf(http.StatusBadRequest, 400, "author id must >= 0")
	}
	return &v1.GetArticlesResp{
		Total: 1,
		Articles: []*v1.Article{
			{
				Title:    "test article: " + req.Title,
				Content:  "test",
				AuthorId: 1,
			},
		},
	}, nil
}

func main() {
	e := mq.NewServer(nil, nil)
	v1.RegisterBlogServiceMQServer(e, &service{})
	e.Start(context.Background())
}
