import Authorization from "./Authorization";

const ErrorCode: { [key: number]: string } = {
    0: "error_null",
    1: "error_server_side",
    2: "error_system_side_test",

    6: "error_client_invalid_request",

    11: "error_username_existed",
    12: "error_email_existed",
    13: "error_register_verification_code_invalid",
    14: "error_invalid_username",
    15: "error_username_too_short",
    16: "error_username_too_long",
    17: "error_invalid_email_format",
    18: "error_invalid_password",
    19: "error_password_too_short",

    21: "error_user_not_exist",
    22: "error_password_not_match",
    23: "error_user_deleted",

    31: "error_username_is_empty",
    32: "error_email_is_empty",
    33: "error_password_is_empty",
    34: "error_new_password_is_empty",

    41: "error_authorization_invalid",
    42: "error_authorization_expired",
    43: "error_not_logged_in",
    44: "error_user_not_match",
    45: "error_user_data_edited",

    51: "error_email_verify_not_exist",
    52: "error_email_verification_code_is_empty",
    53: "error_invalid_email_verification_key"
}

const Host = "/api";

export const Path: { [key: string]: string } = {
    USER: Host + "/user/",
    //GET: request user data
    LOGIN: Host + "/login/",
    //POST: login
    EMAIL: Host + "/email/",
    //GET: verify email change /email/*
    //POST: verify email owner 
    REGISTER: Host + "/register/",
    //GET: verify register /register/*
    //POST register
    USER_USERNAME: Host + "/user/username",
    //PUT: request change username
    USER_EMAIL: Host + "/user/email",
    //PUT: request change email
    //  need verification code
    USER_PASSWORD: Host + "/user/password"
    //PUT: request change password
}

type RawResponseData = {
    success: boolean
    result: string
    error: number
}

export function Post(path: string, postBody?: { [key: string]: string }): Promise<ResponseData> {
    return Request("POST", path, postBody);
}

export function Put(path: string, putBody?: { [key: string]: string }): Promise<ResponseData> {
    return Request("PUT", path, putBody);
}

export function Get(path: string): Promise<ResponseData> {
    return Request("GET", path);
}

function Request(method: string, path: string, body?: { [key: string]: string }): Promise<ResponseData> {
    return new Promise(async (reslove, reject) => {

        const data = {
            method,
            headers: {
                "Content-Type": "application/json",
                "Authorization": Authorization.GetToken()
            },
        }

        if (body !== null) {
            data["body"] = JSON.stringify(body);
        }
        const response = await fetch(path, data);

        if (response.headers.has("authorization")) {
            const authorization = response.headers.get("authorization");
            Authorization.SetToken(authorization);
        }

        if (!response.ok) reject();

        let responseBody = "";
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            responseBody += decoder.decode(value, { stream: true });
        }

        try {
            const bodyJson: RawResponseData = JSON.parse(responseBody);
            const responseData: ResponseData = {
                success: bodyJson.success,
                result: bodyJson.result,
                error: ErrorCode[bodyJson.error]
            }
            reslove(responseData)
        }
        catch (err) {
            reject(err);
        }
    });
}

function GetMe(): Promise<ResponseData> {
    return Get(Path.USER)
}

function Login(email: string, password: string): Promise<ResponseData> {
    return Post(Path.LOGIN, { email, password });
}

function VerifyEmailChange(key: string): Promise<ResponseData> {
    return Get(Path.EMAIL + key);
}

function VerifyEmailOwner(): Promise<ResponseData> {
    return Post(Path.EMAIL);
}

function Register(username: string, email: string, password: string): Promise<ResponseData> {
    return Post(Path.REGISTER, { username, email, password });
}

function VerifyRegister(key: string): Promise<ResponseData> {
    return Get(Path.REGISTER + key);
}

function ChangeUsername(username: string, password: string): Promise<ResponseData> {
    return Put(Path.USER_USERNAME, { username, password });
}

function ChangeEmail(verificationCode: string, newEmail: string): Promise<ResponseData> {
    return Put(Path.USER_EMAIL, { code: verificationCode, email: newEmail });
}

function ChangePassword(currentPassword: string, newPassword: string): Promise<ResponseData> {
    return Put(Path.USER_PASSWORD, { password: currentPassword, new_password: newPassword });
}

export default {
    GetMe,
    Login,
    VerifyEmailChange,
    VerifyEmailOwner,
    Register,
    VerifyRegister,
    ChangeUsername,
    ChangeEmail,
    ChangePassword
}