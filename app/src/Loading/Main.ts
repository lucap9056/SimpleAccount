import { writable, type Writable } from "svelte/store";

class Loading {
    private id: string;
    private parent: Loadings;
    constructor(parent: Loadings, id: string) {
        this.parent = parent;
        this.id = id;
    }

    public Remove(): void {
        this.parent.Remove(this.id);
    }
}

class Loadings {
    public Visible: Writable<boolean>;
    private loadings: { [key: string]: Loading } = {};
    constructor() {
        this.Visible = writable(false);
    }

    public Append(): Loading {
        const id = this.GenerateID();
        const loading = new Loading(this, id);
        this.loadings[id] = loading;
        this.Update();
        return loading;
    }

    public Remove(id: string): void {
        if (this.loadings[id] == null) return;
        delete this.loadings[id];
        this.Update();
    }

    private Update(): void {
        const { loadings, Visible } = this;
        const visible = Object.keys(loadings).length > 0;
        Visible.update(() => visible);
    }

    private GenerateID(): string {
        const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let id = '';
        for (let i = 0; i < 8; i++) {
            id += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return id;
    }
}

const loadings = new Loadings();
export default loadings;
export const visible = loadings.Visible;