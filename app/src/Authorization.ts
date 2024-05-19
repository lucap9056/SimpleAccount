const storage = window.localStorage;
delete window.localStorage;

const TOKEN = "tk";
const TEMP_TOKEN = "_tk";
function SetToken(token: string): void {
    switch (token) {
        case "invalid":
            storage.removeItem(TOKEN);
            storage.removeItem(TEMP_TOKEN);
            break;
        case "invalid_t":
            storage.removeItem(TEMP_TOKEN);
            break;
        default:
            if (/^T /.test(token)) {
                storage.setItem(TEMP_TOKEN, token);
                return;
            }
            if (/\./.test(token)) {
                storage.setItem(TOKEN, token);
            }
    }
}

function GetToken(): string {
    return storage.getItem(TEMP_TOKEN) || storage.getItem(TOKEN) || "";
}

function GetUser(): Promise<User> {
    return new Promise((resolve, reject) => {
        try {
            const token = storage.getItem(TOKEN).split(/\./g);
            const userStr = atob(token[1])
            const user: User = JSON.parse(userStr);
            resolve(user);
        }
        catch {
            reject();
            return;
        }
    });
}

export default {
    SetToken,
    GetToken,
    GetUser
}