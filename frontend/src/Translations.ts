import { writable, type Writable } from "svelte/store";

type Translation = {
    [key: string]: string
}

const _default = {
    "register": "Register",
    "login": "Login",
    "register_username": "Username:",
    "register_email": "Email",
    "register_password": "Password",
    "register_retype_password": "Retype Password",
    "register_email_verification": "An email verification has been sent to your email inbox. Please check your email and confirm to complete the account registration.",
    "register_confirm": "Confirm",
    "register_invalid_email": "Invalid email format.",
    "register_invalid_username": "Username contains invalid characters.",
    "register_invalid_password": "Password contains invalid characters.",
    "register_username_limit": "Available characters: a~zA~Z0~9_-",
    "register_password_limit": "Available characters: a~zA~Z0~9!@#$%^&*_-",
    "register_username_too_short": "Your username must be at least 4 characters long.",
    "register_password_too_short": "Your password must be at least 8 characters long.",
    "register_username_too_long": "Usernames cannot be more than 24 characters.",
    "register_passwords_not_match": "Your passwords do not match.",
    "register_success": "Register success.",

    "login_email": "Email",
    "login_password": "Password",

    "index_info": "User info",
    "index_username": " Name: ",
    "index_email": "Email: ",
    "index_current_password": "Current Password",
    "index_edit_confirm": "Confirm",
    "index_edit_cancel": "Cancel",
    "index_change_username": "Change Username",
    "index_change_email": "Change Email Address",
    "index_change_verify_old_email_address": "Verify Your Old Email Address",
    "index_change_email_send_verification": "Send Verification Code",
    "index_change_new_email_address": "New Email Address",
    "index_change_email_again_count_down": "waiting:",
    "index_change_password": "Change Password",
    "index_change_password_cant_equal": "Your new password can not be equal to your current password.",
    "index_change_current_password": "Current Password",
    "index_change_new_password": "New Password",
    "index_change_retype_password": "Retype Password",
    "index_change_success": "Success",
    "index_change_check_new_email_address": "Please go to your new email inbox and click on the link to apply the changes.",
    "index_change_verification_email_owner": "Please check your email for the verification code.",

    "error_null": "",
    "error_server_side": "Server side error",
    "error_server_side_test": "Server side test",
    "error_client_invalid_request": "Invalid request",
    "error_username_exists": "Username already exists",
    "error_email_existed": "User email already exists",
    "error_register_verification_code_invalid": "Verification code invalid",
    "error_invalid_username": "Invalid username",
    "error_username_too_short": "Username too short",
    "error_username_too_long": "Username too long",
    "error_invalid_email_format": "Invalid email address format",
    "error_invalid_password": "Invalid password",
    "error_password_too_short": "Password too short",
    "error_user_not_exist": "User does not exist",
    "error_password_not_match": "Passwords do not match",
    "error_username_is_empty": "Username is empty",
    "error_email_is_empty": "Email Address is empty",
    "error_password_is_empty": "Password is empty",
    "error_new_password_is_empty": "New password is empty",
    "error_authorization_invalid": "Authorization invalid",
    "error_authorization_expired": "Authorization expired",
    "error_not_logged_in": "Not logged in",
    "error_user_not_match": "User does not match",
    "error_email_verify_not_exist": "Email address verification does not exist",
    "error_email_verification_code_is_empty": "Email verification code is empty",
    "error_invalid_email_verification_key": "Invalid verification key"
};

let translation: Translation = _default;
const Translations: Writable<Translation> = writable(translation);
Translations.subscribe((l) => {
    translation = l;
});

fetch(`/translations/${navigator.language}.json`).then(async (res) => {
    if (!res.ok) return;

    let body = "";
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        body += decoder.decode(value, { stream: true });
    }

    const Translation: Translation = Object.assign(_default, JSON.parse(body));
    Translations.update(() => Translation);
});

export default Translations;

export function TranslationsGet(id: string) {
    return translation[id] || id;
}