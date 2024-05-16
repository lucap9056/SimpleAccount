package simple_account_http_post

import (
	"simple_account/app/Email/TimedKeys"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func EmailOwner(context *Message.Context) (int, error) {
	if context.Author.User == nil {
		return Error.NOT_LOGGED_IN, nil
	}

	author := context.Author.User

	shortKey, err := TimedKeys.GenerateShortKey()
	if err != nil {
		return Error.SYSTEM, nil
	}
	key := author.Username + "-" + shortKey
	err = context.Email.TimedKeys.Add("email_owner", key, *author)
	if err != nil {
		return Error.SYSTEM, err
	}

	language := context.Author.Language
	err = context.Email.Templates.EmailOwner.SendVerify(language, author.Email, shortKey)
	if err != nil {
		return Error.SYSTEM, err
	}

	return Error.NULL, nil
}
