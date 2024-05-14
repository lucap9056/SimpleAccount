export class AlertManager {
    public alerts: Alert[] = [];
    private listen: (alerts: Alert[]) => void;

    public SetListen(listen: (alerts: Alert[]) => void) {
        this.listen = listen;
    }

    public GetNotifys() {
        return this.alerts;
    }

    public Add(message: string, type: number, asyncFunc?: (callback?: () => void) => Promise<void>, removeMessage?: string): Alert {
        const { alerts, listen } = this;
        const alert = new Alert(this, message, type, removeMessage, asyncFunc);
        const id = alert.GetId();

        if (alerts.find((n) => n.GetId() == id)) {
            throw "Alert Exists";
        }
        alerts.push(alert);

        if (listen) listen(alerts);

        if (asyncFunc) {
            switch (type) {
                case Alert.Type.Normal:
                case Alert.Type.Error:
                    asyncFunc(alert.Remove.bind(alert)).catch((err) => {
                        alert.Remove();
                        alertManager.Add(err, Alert.Type.Error);
                    });
                    break;
            }
        }
        return alert;
    }

    public Remove(id: string): void {
        const { alerts, listen } = this;
        const index = alerts.findIndex((n) => n.GetId() == id);
        if (index == -1) return;
        alerts.splice(index, 1);
        if (listen) listen(alerts);
    }
}

export const alertManager = new AlertManager();

export class Alert {
    public static Type = class {
        public static Normal: number = 0;
        public static Alert: number = 1;
        public static Error: number = 2;
    }
    private manager: AlertManager;
    private message: string;
    private removeMessage: string;
    private type: number;
    private id: string;
    private func: () => void;

    constructor(manager: AlertManager, message: string, type: number = 0, removeMessage: string = "", alertFunc?: () => void) {
        this.manager = manager;
        this.id = this.CreateId();
        this.message = message;
        this.removeMessage = removeMessage;
        this.type = type;

        switch (type) {
            case Alert.Type.Alert:
                this.func = alertFunc;
                break;
            case Alert.Type.Normal:
                setTimeout(this.Remove.bind(this), 3000);
                break;
            case Alert.Type.Error:
                setTimeout(this.Remove.bind(this), 5000);
                break;
        }
    }

    private CreateId(): string {
        let buf = new Uint16Array(8);
        crypto.getRandomValues(buf);
        return ([buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7]]).map(function (num) {
            return (num < 0x10000 ? '0' : '') + num.toString(16);
        }).join('-');
    }

    public GetId(): string {
        return this.id;
    }

    public GetType(): number {
        return this.type;
    }

    public GetMessage(): string {
        return this.message;
    }

    public GetRemoveMessage(): string {
        return this.removeMessage;
    }

    public Remove(): void {
        this.manager.Remove(this.id);

        if (this.type === Alert.Type.Alert) {
            if (this.func) this.func();
        }
    }
}

export default {
    AlertManager,
    Alert,
    alertManager
}