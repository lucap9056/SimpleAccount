import { writable, type Writable } from "svelte/store";
import { Get, Path } from "./API";
import { Alert, alertManager } from "./Alert/Struct";
import Loading from "./Loading/Main"
import Translations from "./Translations";
import Status from "./Status";

let login: boolean = false;
Status.login.subscribe((l) => {
    login = l;
});

export class Hash {
    private dirs: string[];
    constructor(hash: string) {
        this.dirs = hash.replace(/^#/, "").split("/");
    }

    public Shift(): string {
        if (this.dirs.length > 0) {
            return this.dirs.shift();
        }
        return "";
    }
}

export const Routes: { [key: string]: string } = {
    INDEX: "INDEX",
    LOGIN: "LOGIN",
    REGISTER: "REGISTER",
};

export const router = new class {
    public route: Writable<Route>;
    private hash: Hash;
    constructor() {
        const routes = this.CreateRoutes();
        this.route = writable(routes);

        this.hash = new Hash(location.hash);
        window.addEventListener("hashchange", this.HashChange.bind(this));
        this.HashChange();
    }

    public Set(...routes: string[]): void {
        if (routes.length < 1) return;
        let hash = "#" + routes[0].toLowerCase();

        for (let i = 1; i < routes.length; i++) {
            hash += "/" + routes;
        }
        location.hash = hash;
    }

    private CreateRoutes(): Route {
        const routes: Route = {};
        for (const route of Object.keys(Routes)) {
            routes[route] = false;
        }
        return routes;
    }

    private async HashChange() {
        const { CreateRoutes, Set } = this;
        const hash = new Hash(location.hash);
        this.hash = hash;

        const routes = CreateRoutes();

        const head = hash.Shift().toUpperCase();

        let exist = false;
        for (const r of Object.values(Routes)) {
            if (head == r) {
                exist = true;
                routes[head] = true;
            }
        }

        if (!exist) {
            if (login) {
                Set(Routes.INDEX);
            } else {
                Set(Routes.LOGIN);
            }
            return;
        }

        this.route.update(() => routes);

        switch (head) {
            case Routes.REGISTER:
                const id = hash.Shift();
                if (id === "") break;
                const loading = Loading.Append();
                try {
                    const res = await Get(Path.REGISTER + id);
                    if (res.success) {
                        alertManager.Add(
                            Translations.Get("register_success"),
                            Alert.Type.Alert,
                            null,
                            Translations.Get("register_confirm"),
                        );
                    } else {
                        alertManager.Add(Translations.Get(res.error), Alert.Type.Error);
                    }
                    Set(Routes.LOGIN);
                }
                catch { }
                loading.Remove();

                break;
        }
    }

    public Hash() {
        return this.hash;
    }
}

export const route = router.route;

export default {
    router,
    route,
    Routes
}