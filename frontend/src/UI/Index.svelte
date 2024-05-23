<script lang="ts">
    import Authorization from "../Authorization";
    import { Routes, router } from "../Router";
    import Status from "../Status";
    import Translations from "../Translations";
    import EmailForm from "../EditUI/Email.svelte";
    import UsernameForm from "../EditUI/Username.svelte";
    import PasswordForm from "../EditUI/Password.svelte";

    let user: User;
    Authorization.GetUser().then((usr) => {
        user = usr;
    });

    function logout() {
        Authorization.SetToken("invalid");
        Authorization.SetToken("invalid_t");
        Status.login.update(() => false);
        router.Set(Routes.LOGIN);
    }

    let editUsername = false;
    function EditUsername() {
        editUsername = true;
    }

    let editEmail = false;
    function EditEmail() {
        editEmail = true;
    }

    let editPassword = false;
    function EditPassword() {
        editPassword = true;
    }

    function EditCancel() {
        Authorization.GetUser().then((usr) => {
            user = usr;
        });
        editUsername = false;
        editEmail = false;
        editPassword = false;
    }
</script>

<div class="container">
    {#if user}
        {#if !(editUsername || editEmail || editPassword)}
            <div class="user-info">
                <h2>{$Translations.index_info}</h2>
                <div
                    class="field"
                    data-label={$Translations.index_username}
                >
                    <input type="text" bind:value={user.name} readonly />
                    <button class="field_edit" on:click={EditUsername}>
                        <ion-icon name="settings-outline"></ion-icon>
                    </button>
                </div>
                <div class="field" data-label={$Translations.index_email}>
                    <input type="text" bind:value={user.email} readonly />
                    <button class="field_edit" on:click={EditEmail}>
                        <ion-icon name="settings-outline"></ion-icon>
                    </button>
                </div>

                <div class="field">
                    <button class="password_edit btn" on:click={EditPassword}>
                        {$Translations.index_change_password}
                    </button>

                    <button id="logout" class="btn right" on:click={logout}>
                        <ion-icon name="log-out-outline"></ion-icon>
                    </button>
                </div>
            </div>
        {:else}
            {#if editUsername}
                <UsernameForm Cancel={EditCancel} />
            {/if}
            {#if editEmail}
                <EmailForm Cancel={EditCancel} />
            {/if}
            {#if editPassword}
                <PasswordForm Cancel={EditCancel} />
            {/if}
        {/if}
    {/if}
</div>

<style>
    .container {
        max-width: 400px;
        margin: 0 auto;
        padding: 20px;
        border: 1px solid #ccc;
        border-radius: 5px;
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    }

    .user-info {
        margin-bottom: 20px;
        padding-bottom: 16px;
    }

    .user-info h2 {
        margin-top: 0;
        font-size: 24px;
        color: #333;
    }

    .field {
        position: relative;
        margin: 0;
        font-size: 16px;
        color: #666;
        margin-bottom: 0.5em;
    }

    .field::before {
        content: attr(data-label);
        position: relative;
        text-align: right;
        line-height: 32px;
        float: left;
        width: 100px;
    }

    input {
        width: 240px;
        outline: none;
        float: left;
        background-color: transparent;
        border-color: transparent;
    }

    .btn {
        padding: 8px 12px 5px 12px;
        margin-right: 10px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        font-size: 16px;
        color: #fff;
        background-color: #007bff;
        transition: background-color 0.3s;
    }

    .field_edit {
        height: 32px;
        font-size: 10px;
        border-radius: 3px;
        transition: background-color 0.3s;
    }

    .field_edit:hover {
        background-color: #ccc;
    }

    .field_edit ion-icon {
        font-size: 24px;
    }

    .password_edit {
        width: auto;
        height: 32px;
        line-height: 32px;
        padding: 0;
        margin-left: 100px;
        float: left;
        padding: 0 20px;
    }

    .right {
        float: right;
    }

    .btn:hover {
        background-color: #0056b3;
    }

    #logout {
        height: 32px;
        padding: 4px 16px;
    }
    #logout ion-icon {
        font-size: 24px;
    }
</style>
