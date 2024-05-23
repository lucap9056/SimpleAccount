import { writable, type Writable } from "svelte/store";
import API from "./API";

export const login: Writable<boolean> = writable(false);

export function Login(): Promise<boolean> {
    return new Promise((reslove) => {
        API.GetMe().then((res) => {
            login.update(() => true);
            reslove(true);
        }).catch(() => {
            API.GetMe().then(() => {
                login.update(() => true);
                reslove(true);
            }).catch(() => {
                reslove(false);
            });
        });
    });
}

export default {
    Login,
    login
}