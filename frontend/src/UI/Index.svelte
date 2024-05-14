<script lang="ts">
    import Authorization from "../Authorization";
    import { Routes, router } from "../Router";
    import Status from "../Status";
    import Translations from "../Translations";

    let edit: boolean = false;
    let user: User;
    Authorization.GetUser().then((usr) => {
        user = usr;
    });

    function save() {
        //edit = false;
    }

    function editProfile() {
        //edit = true;
    }

    function logout() {
        Authorization.SetToken("invalid");
        Authorization.SetToken("invalid_t");
        Status.login.update(() => false);
        router.Set(Routes.LOGIN);
    }
</script>

<div class="container">
    {#if user}
        <div class="user-info">
            <h2>{Translations.Get("index_info")}</h2>
            <div class="input" data-label={Translations.Get("index_username")}>
                <input type="text" bind:value={user.name} readonly={!edit} />
            </div>
            <div class="input" data-label={Translations.Get("index_email")}>
                <input type="text" bind:value={user.email} readonly={!edit} />
            </div>
        </div>
        <div class="options">
            {#if edit}
                <button class="right btn" on:click={save}>
                    <ion-icon name="save-outline"></ion-icon>
                </button>
            {:else}
                <button class="left btn">
                    <ion-icon name="lock-closed-outline"></ion-icon>
                </button>
                <button class="right btn" on:click={editProfile}>
                    <ion-icon name="settings-outline"></ion-icon>
                </button>
                <button class="right btn" on:click={logout}>
                    <ion-icon name="log-out-outline"></ion-icon>
                </button>
            {/if}
        </div>
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
    }

    .user-info h2 {
        margin-top: 0;
        font-size: 24px;
        color: #333;
    }

    .input {
        position: relative;
        margin: 0;
        font-size: 16px;
        color: #666;
        margin-bottom: 0.5em;
    }

    .input::before {
        content: attr(data-label);
        position: absolute;
        right: calc(100% - 110px);
        text-align: right;
        line-height: 32px;
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
    .btn.left {
        float: left;
    }
    .btn.right {
        float: right;
    }

    .btn:hover {
        background-color: #0056b3;
    }

    input {
        width: 240px;
        margin-left: 60px;
    }

    input:read-only {
        outline: none;
    }

    .options {
        height: 42px;
    }

    ion-icon {
        font-size: 24px;
    }
</style>
