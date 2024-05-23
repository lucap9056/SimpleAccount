<script lang="ts">
    export let name: string;
    export let label: string = null;
    export let value: string;
    export let placeholder: string = "";
    export let password: boolean = false;
    export let hint: string = "";
    export let validate: (value: string) => string = null;

    let rawValue = value.toString();
    let failed = "";

    $: {
        if (validate !== null) {
            failed = validate(rawValue);
        }
        if (rawValue !== "" && failed == "") {
            value = rawValue;
        } else {
            value = "";
        }
    }
</script>

<div class="form-group" data-hint={rawValue !== "" && failed !== ""}>
    <label for={name}>{label}</label>
    <div class="input">
        {#if password}
            <input
                type="password"
                id={name}
                bind:value={rawValue}
                required
                {placeholder}
            />
        {:else}
            <input
                type="text"
                id={name}
                bind:value={rawValue}
                required
                {placeholder}
            />
        {/if}
        {#if rawValue !== "" && failed !== ""}
            <div class="input_alert">{failed}</div>
        {/if}
        {#if hint != ""}
            <div class="hint">
                {hint}
            </div>
        {/if}
    </div>
</div>

<style>
    .form-group {
        margin-top: 5px;
        margin-bottom: 0px;
        transition-duration: 25ms;
    }

    .form-group[data-hint="true"] {
        margin-bottom: 40px;
    }

    label {
        display: block;
        margin-bottom: 5px;
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

    input[type="text"],
    input[type="password"] {
        width: 100%;
        padding: 10px;
        border-radius: 5px;
        border: 1px solid #ccc;
        box-sizing: border-box;
    }

    .hint {
        font-size: 12px;
        color: #999;
        user-select: none;
        text-align: left;
        margin: 2px 5px;
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
