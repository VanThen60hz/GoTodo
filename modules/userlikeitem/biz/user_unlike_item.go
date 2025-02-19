package biz

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"GoTodo/pubsub"
	"context"
	"log"
)

type UserUnLikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

//type DecrementLikeCountStore interface {
//	DecreaseLikeCount(ctx context.Context, id int) error
//}

type userUnLikeItemBiz struct {
	store UserUnLikeItemStore
	//itemStore DecrementLikeCountStore
	ps pubsub.PubSub
}

func NewUserUnLikeItemBiz(
	store UserUnLikeItemStore,
	//itemStore DecrementLikeCountStore,
	ps pubsub.PubSub,
) *userUnLikeItemBiz {
	return &userUnLikeItemBiz{
		store: store,
		//itemStore: itemStore,
		ps: ps,
	}
}

func (biz *userUnLikeItemBiz) UnLikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)

	if err == common.RecordNotFound {
		return model.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.ps.Publish(ctx, common.TopicUserUnlikeItem, pubsub.NewMessage(&model.Like{UserId: userId, ItemId: itemId})); err != nil {
		log.Println(err)
	}

	//job := asyncjob.NewJob(func(ctx context.Context) error {
	//	if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
	//		return err
	//	}
	//
	//	return nil
	//})
	//
	//if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	//	log.Println(err)
	//}

	return nil
}
