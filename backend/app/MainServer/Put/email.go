package simple_account_http_put

import (
	"simple_account/app/Email/TimedKeys"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
)

func Email(context *Message.Context) (int, error) {
	if context.Author.User == nil {
		return Error.NOT_LOGGED_IN, nil
	}

	author, errCode, err := context.Database.GetUser(context.Author.User.Id)
	if errCode != Error.NULL {
		return errCode, err
	}

	email := context.Email

	userData, errCode, err := context.ResquestBody()
	if errCode != Error.NULL {
		return errCode, err
	}

	if userData.VerificationCode == "" {
		return Error.EMAIL_VERIFICATION_CODE_IS_EMPTY, nil
	}

	if userData.Email == "" {
		return Error.EMAIL_IS_EMPTY, nil
	}

	verifyKey := author.Username + "-" + userData.VerificationCode
	user := email.TimedKeys.Verify("email_owner", verifyKey)
	if user == nil {
		return Error.INVALID_EMAIL_VERIFICATION_KEY, nil
	}

	if author.Email != user.Email || author.Id != user.Id {
		return Error.USER_NOT_MATCH, nil
	}

	key, err := TimedKeys.GenerateKey()
	if err != nil {
		return Error.SYSTEM, err
	}

	language := context.Author.Language
	err = email.Templates.ChangeEmail.SendVerify(language, userData.Email, key)
	if err != nil {
		return Error.SYSTEM, err
	}

	updatedUser := author
	updatedUser.Email = userData.Email

	err = email.TimedKeys.Add("change_email", key, updatedUser)
	if err != nil {
		return Error.SYSTEM, err
	}

	email.TimedKeys.Del("email_owner", userData.VerificationCode)

	return Error.NULL, nil
}
