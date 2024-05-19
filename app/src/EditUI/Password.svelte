<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import Translations from "../Translations";
    import Validate from "../ValidateInput";
    import API from "../API";
    import loadings from "../Loading/Main";
    import { Alert, alertManager } from "../Alert/Struct";
    const dispatch = createEventDispatcher();
    function Cancel() {
        dispatch("cancel");
    }

    let currentPassword = "";
    let newPassword = "";
    let retypePassword = "";

    let newPasswordFailed = "";
    let retypePasswordFailed = "";
    $: {
        newPasswordFailed = Validate.Password(newPassword, currentPassword);
        retypePasswordFailed = Validate.RetypePassword(
            newPassword,
            retypePassword,
        );
    }

    function handleSubmit() {
        if (newPasswordFailed != "" || retypePasswordFailed != "") {
            return;
        }

        const loading = loadings.Append();

        API.ChangePassword(currentPassword, newPassword).then((res) => {
            if (res.success) {
                Cancel();
                alertManager.Add(
                    $Translations.index_change_success,
                    Alert.Type.Normal,
                );
            } else {
                alertManager.Add($Translations[res.error], Alert.Type.Error);
            }
            loading.Remove();
        });
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <h2>{$Translations.index_change_password}</h2>

    <div class="form-group">
        <label for="current_password">
            {$Translations.index_change_current_password}
        </label>
        <div class="input" data-type="current_password">
            <input
                type="password"
                id="current_password"
                bind:value={currentPassword}
                required
            />
        </div>
    </div>

    <div class="form-group">
        <label for="new_password">
            {$Translations.index_change_new_password}
        </label>
        <div class="input" data-type="new_password">
            <input
                type="password"
                id="new_password"
                bind:value={newPassword}
                required
            />
            {#if newPassword != "" && newPasswordFailed !== ""}
                <div class="input_alert">{newPasswordFailed}</div>
            {/if}

            <div class="limit">
                {$Translations.register_password_limit}
            </div>
        </div>
    </div>

    <div class="form-group">
        <label for="retype_password">
            {$Translations.index_change_retype_password}
        </label>
        <div class="input" data-type="retype_password">
            <input
                type="password"
                id="retype_password"
                bind:value={retypePassword}
                required
            />

            {#if retypePassword != "" && retypePasswordFailed !== ""}
                <div class="input_alert">{retypePasswordFailed}</div>
            {/if}
        </div>
    </div>

    <div class="form-group options">
        <button type="button" on:click={Cancel}>
            {$Translations.index_edit_cancel}
        </button>
        <button type="submit">
            {$Translations.index_edit_confirm}
        </button>
    </div>
</form>

<style>
    h2 {
        text-align: center;
    }

    .form-group {
        margin-bottom: 20px;
    }

    .options {
        display: flex;
        justify-content: flex-end;
    }

    label {
        display: block;
        margin-bottom: 5px;
    }

    .input {
        position: relative;
    }

    input {
        position: relative;
    }

    input {
        width: 100%;
        padding: 10px;
        border-radius: 5px;
        border: 1px solid #ccc;
        box-sizing: border-box;
    }

    button {
        width: 80px;
        padding: 10px;
        background-color: #007bff;
        color: #fff;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        margin-left: 20px;
    }

    button:hover {
        background-color: #0056b3;
    }

    .limit {
        font-size: 12px;
        color: #999;
        user-select: none;
        text-align: left;
    }

    .input_alert {
        position: absolute;
        top: 100%;
        white-space: nowrap;
        height: 32px;
        line-height: 32px;
        padding: 0 10px;
        background-color: #c46;
        border-radius: 5px;
        color: white;
        margin: 4px;
    }
</style>
