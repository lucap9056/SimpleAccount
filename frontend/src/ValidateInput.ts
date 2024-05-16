import Translations from "./Translations";

function Username(value: string): string {
    if (!/^[a-zA-Z0-9_-]+$/.test(value)) {
        return Translations.Get("register_invalid_username");
    }

    if (value.length < 4) {
        return Translations.Get("register_username_too_short");
    }

    if (value.length > 24) {
        return Translations.Get("register_username_too_length");
    }

    return "";
}

function Email(value: string): string {
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        return Translations.Get("register_invalid_email");
    }
    return "";
}

function Password(value: string, old?: string): string {
    if (value.length < 8) {
        return Translations.Get("register_password_too_short");
    }

    if (!/^[A-Za-z0-9!@#$%^&*_\-]+$/.test(value)) {
        return Translations.Get("register_invalid_password");
    }

    if (old && value === old) {
        return Translations.Get("index_change_password_cant_equal")
    } 
    return "";
}

function RetypePassword(o: string, r: string) {
    if (o === r) return "";
    return Translations.Get("register_passwords_not_match");
}

export default {
    Username,
    Email,
    Password,
    RetypePassword
}