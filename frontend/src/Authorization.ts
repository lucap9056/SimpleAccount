const storage = window.localStorage;
delete window.localStorage;

const TOKEN = "tk";
const TEMP_TOKEN = "_tk";
function SetToken(token: string): void {
    switch (token) {
        case "invalid":
            storage.removeItem(TOKEN);
            break;
        case "invalid_t":
            storage.removeItem(TEMP_TOKEN);
            break;
        default:
            if (/^Bearer /.test(token)) {
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
        const token = storage.getItem(TOKEN);
        try {
            const playloadStr = atob(token.replace(/\..*/, ''));
            const playload: Playload = JSON.parse(playloadStr);
            resolve(playload.user);
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