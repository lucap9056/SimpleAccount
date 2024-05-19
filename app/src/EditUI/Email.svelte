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

    let verification = "";
    let newEmail = "";

    let newEmailFalied = "";
    $: {
        newEmailFalied = Validate.Email(newEmail);
    }

    let againCountDown = 0;

    function GetVerificationCode() {
        if (againCountDown > 0) {
            return;
        }
        againCountDown = 120;
        const countDown = setInterval(() => {
            if (againCountDown == 0) {
                clearInterval(countDown);
            } else {
                againCountDown--;
            }
        }, 1000);

        const loading = loadings.Append();
        API.VerifyEmailOwner().then((res) => {
            if (res.success) {
                alertManager.Add(
                    $Translations.index_change_verification_email_owner,
                    Alert.Type.Alert,
                );
            } else {
                alertManager.Add($Translations[res.error], Alert.Type.Error);
            }
            loading.Remove();
        });
    }

    function handleSubmit() {
        if (newEmailFalied != "") {
            return;
        }

        const loading = loadings.Append();

        API.ChangeEmail(verification, newEmail).then((res) => {
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
    <h2>{$Translations.index_change_email}</h2>

    <div class="form-group">
        <label for="verification">
            {$Translations.index_change_verify_old_email_address}
        </label>
        <div class="input" data-type="verification_code">
            <input
                type="text"
                id="verification"
                bind:value={verification}
                required
            /><button type="button" on:click={GetVerificationCode}>
                {#if againCountDown == 0}
                    {$Translations.index_change_email_send_verification}
                {:else}
                    {$Translations.index_change_email_again_count_down}{againCountDown}s
                {/if}
            </button>
        </div>
    </div>

    <div class="form-group">
        <label for="new_email">
            {$Translations.index_change_new_email_address}
        </label>

        <div class="input">
            <input type="text" id="new_email" bind:value={newEmail} required />

            {#if newEmail != "" && newEmailFalied !== ""}
                <div class="input_alert">{newEmailFalied}</div>
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
        padding: 10px;
        border-radius: 5px;
        border: 1px solid #ccc;
        box-sizing: border-box;
    }

    #verification {
        width: 80px;
        border-radius: 5px 0 0 5px;
    }

    #new_email {
        width: 100%;
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

    #verification + button {
        margin-left: 0;
        border-radius: 0 5px 5px 0;
        width: auto;
    }

    button:hover {
        background-color: #0056b3;
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
