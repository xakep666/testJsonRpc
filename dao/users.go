package dao

import "testJsonRpc/model"

type UserEditArgs struct {
	Login string
	NewData model.User
}

// this type is for rpc package
type Users int

func (Users) Register(login string, ok *bool) error {
	err:= dbImpl.Register(login)
	*ok=err==nil
	return err
}

func (Users) GetByLogin(login string, resp *model.User) error {
	user, err:= dbImpl.GetByLogin(login)
	*resp=user
	return err
}

func (Users) Edit(args UserEditArgs, ok *bool) error {
	err:= dbImpl.Edit(args.Login, args.NewData)
	*ok=err==nil
	return err
}
