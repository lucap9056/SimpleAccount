const _default = {
    register: "Register",
    login: "Login",

    register_username: "Username:",
    register_email: "Email:",
    register_password: "Password:",
    register_retype_password: "Retype Password:",
    register_email_verification: "An email verification has been sent to your email inbox. Please check your email and confirm to complete the account registration.",
    register_confirm: "Confirm",
    register_invalid_email: "Invalid email format.",
    register_invalid_username: "Username contains invalid characters.",
    register_invalid_password: "Password contains invalid characters.",
    register_username_limit: "Available characters: a~zA~Z0~9_-",
    register_password_limit: "Available characters: a~zA~Z0~9!@#$%^&*_-",
    register_username_too_short: "Your username must be at least 4 characters long.",
    register_password_too_short: "Your password must be at least 8 characters long.",
    register_username_too_length: "Usernames cannot be more than 24 characters.",
    register_passwords_not_match: "Your passwords do not match.",
    register_success: "Register success.",

    login_username: "Name",
    login_password: "Password",

    index_info: "User info",
    index_username: " Name: ",
    index_email: "Email: ",

    error_null: "",
    error_server_side: "Server side error",
    error_server_side_test: "Server side test",
    error_client_invalid: "Invalid request",
    error_register_username_existed: "Username existed",
    error_register_email_existed: "User email existed",
    error_register_verification_code_invalid: "Verification code invalid"
}

const translations: { [key: string]: { [key: string]: string } } = {
    "zh-TW": {
        register: "註冊",
        login: "登入",

        register_username: "使用者名稱：",
        register_email: "電子信箱：",
        register_password: "密碼：",
        register_email_verification: "驗證email已發送至您的電子信箱，請前往確認完成帳戶註冊",
        register_confirm: "確認",
        register_invalid_email: "電子信箱格式無效",
        register_invalid_username: "使用者名稱中包含無效字元",
        register_invalid_password: "密碼中包含無效字元",
        register_username_limit: "可用字元：a~zA~Z0~9_-",
        register_password_limit: "可用字元：a~zA~Z0~9!@#$%^&*_-",
        register_username_too_short: "使用者名稱最短長度為4個字元",
        register_password_too_short: "密碼最短長度為8個字元",
        register_username_too_length: "使用者名稱不可超過24個字元",
        register_password_too_length: "密碼不可超過32個字元",
        register_passwords_not_match: "密碼不相同",
        register_success: "註冊成功",

        login_username: "使用者名稱",
        login_password: "密碼",

        index_info: "使用者資訊",
        index_username: "名稱：",
        index_email: "電子信箱：",

        error_server_side: "伺服端錯誤",
        error_server_side_test: "伺服端測試",
        error_client_invalid: "無效請求",
        error_register_username_existed: "使用者名稱已存在",
        error_register_email_existed: "使用者電子信箱已存在",
        error_register_verification_code_invalid: "驗證碼無效"
    }
}

function Get(id: string): string {
    return translations[navigator.language][id] || _default[id] || id;
}

export default {
    Get
}