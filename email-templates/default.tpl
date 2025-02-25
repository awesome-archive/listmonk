<!doctype html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1" />
        <base target="_blank">

        <style>
            body {
                background-color: #F0F1F3;
                font-family: 'Helvetica Neue', 'Segoe UI', Helvetica, sans-serif;
                font-size: 15px;
                line-height: 26px;
                margin: 0;
                color: #444;
            }

            .wrap {
                background-color: #fff;
                padding: 30px;
                max-width: 525px;
                margin: 0 auto;
                border-radius: 5px;
            }

            .footer {
                text-align: center;
                font-size: 12px;
                color: #888;
            }
                .footer a {
                    color: #888;
                }

            .gutter {
                padding: 30px;
            }

            img {
                max-width: 100%;
            }

            a {
                color: #7f2aff;
            }
                a:hover {
                    color: #111;
                }
            @media screen and (max-width: 600px) {
                .wrap {
                    max-width: auto;
                }
                .gutter {
                    padding: 10px;
                }
            }
        </style>
    </head>
<body style="background-color: #F0F1F3;">
    <div class="gutter">&nbsp;</div>
    <div class="wrap">
        {{ template "content" . }}
    </div>
    
    <div class="footer">
        <p>Don't want to receive these e-mails? <a href="{{ .UnsubscribeURL }}">Unsubscribe</a></p>
        <p>Powered by <a href="https://listmonk.app" target="_blank">listmonk</a></p>
    </div>
    <div class="gutter">&nbsp;{{ TrackView }}</div>
</body>
</html>
