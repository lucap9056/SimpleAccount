<script lang="ts">
    import { onMount } from "svelte";
    import { alertManager, Alert } from "./Struct";
    let container: HTMLDivElement;
    let alerts: Alert[] = alertManager.GetNotifys();
    onMount(() => {
        alertManager.SetListen(Listen);
    });

    function Listen(a: Alert[]) {
        alerts = a;
    }
</script>

{#if alerts.length > 0}
    <div id="alerts_container" bind:this={container}>
        <div class="alert" data-type={alerts[0].GetType()}>
            <div class="alert_message">
                {alerts[0].GetMessage()}
            </div>
            {#if alerts[0].GetType() == Alert.Type.Alert}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div
                    class="remove_alert"
                    on:click={alerts[0].Remove.bind(alerts[0])}
                >
                    {alerts[0].GetRemoveMessage()}
                </div>
            {/if}
        </div>
    </div>
{/if}

<style type="text/scss" lang="scss">
    $alert_normal_color: #fff;
    $alert_alert_color: #fff;
    $alert_error_color: #ee6e67;
    $alert_text_color: #222;
    $alert_border: 1px solid #ccc;
    $alert_shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    $alert_remove_color: #007bff;
    $alert_remove_hover_color: #0056b3;
    $alert_remove_text_color: #fff;

    #alerts_container {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        display: flex;
        justify-content: center;
        align-items: center;

        .alert {
            border: $alert_border;
            border-radius: 4px;
            padding: 10px 20px;
            margin-bottom: 10px;
            box-shadow: $alert_shadow;
            color: $alert_text_color;

            opacity: 0;
            animation: fadeIn 0.3s forwards;
            animation-play-state: running;

            &[data-type="0"] {
                background-color: $alert_normal_color;
            }

            &[data-type="1"] {
                background-color: $alert_alert_color;
            }

            &[data-type="2"] {
                background-color: $alert_error_color;
            }

            .alert_message {
                margin-bottom: 5px;
            }
        }

        .remove_alert {
            cursor: pointer;
            padding: 8px 16px;
            background-color: $alert_remove_color;
            color: $alert_remove_text_color;
            border-radius: 4px;
            font-size: 14px;
            transition: background-color 3s;
            margin: 0 20px;

            &:hover {
                background-color: $alert_remove_hover_color;
            }
        }
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
</style>
