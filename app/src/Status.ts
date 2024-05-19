import { writable, type Writable } from "svelte/store";
import API from "./API";

export const login: Writable<boolean> = writable(false);

export function Login(): Promise<boolean> {
    return new Promise((reslove) => {
        API.GetMe().then((res) => {
            if (res.success) {
                login.update(() => true);
                reslove(true);
            } else {
                API.GetMe().then((resp) => {
                    login.update(() => resp.success);
                    reslove(resp.success);
                });
            }
        });
    });
}

export default {
    Login,
    login
}