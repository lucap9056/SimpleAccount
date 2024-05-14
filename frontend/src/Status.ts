import { writable, type Writable } from "svelte/store";
import { Get, Path } from "./API";

export const login: Writable<boolean> = writable(false);

export function Login(): Promise<boolean> {
    return new Promise((reslove) => {
        Get(Path.USER).then((res) => {
            const loggedIn = res.success;
            login.update(() => loggedIn);
            reslove(loggedIn);
        });
    });
}

export default {
    Login,
    login
}