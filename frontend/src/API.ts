import Authorization from "./Authorization";

const ErrorCode: { [key: number]: string } = {
    0: "error_null",
    1: "error_server_side",
    2: "error_system_side_test",
    6: "error_client_invalid_request",
    11: "error_register_username_existed",
    12: "error_register_email_existed",
    13: "error_register_verification_code_invalid",
}

export const Path: { [key: string]: string } = {
    USER: "/api/user/",
    LOGIN: "/api/login/",
    REGISTER: "/api/register/"
}

type RawResponseData = {
    success: boolean
    result: string
    error: number
}

export function Post(path: string, postBody: { [key: string]: string }): Promise<ResponseData> {
    return new Promise(async (resolve, reject) => {
        const response = await fetch(path, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": Authorization.GetToken()
            },
            body: JSON.stringify(postBody),
        });

        if (response.headers.has("authorization")) {
            const authorization = response.headers.get("authorization");
            Authorization.SetToken(authorization);
        }

        if (!response.ok) reject();

        let body = "";
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            body += decoder.decode(value, { stream: true });
        }

        try {
            const bodyJson: RawResponseData = JSON.parse(body);
            const responseData: ResponseData = {
                success: bodyJson.success,
                result: bodyJson.result,
                error: ErrorCode[bodyJson.error]
            }
            resolve(responseData)
        }
        catch (err) {
            console.log(err);
            console.log(response.body);
            reject();
        }
    });

}

export async function Get(path: string): Promise<ResponseData> {
    return new Promise(async (resolve, reject) => {
        const response = await fetch(path, {
            method: "GET",
            headers: {
                "Authorization": Authorization.GetToken()
            }
        });

        if (response.headers.has("authorization")) {
            const authorization = response.headers.get("authorization");
            Authorization.SetToken(authorization);
        }

        if (!response.ok) reject();

        let body = "";
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            body += decoder.decode(value, { stream: true });
        }

        try {
            const bodyJson: RawResponseData = JSON.parse(body);
            const responseData: ResponseData = {
                success: bodyJson.success,
                result: bodyJson.result,
                error: ErrorCode[bodyJson.error]
            }
            resolve(responseData)
        }
        catch {
            console.log(response.body);
            reject();
        }
    });
}

export default {
    Post,
    Get,
    Path
}