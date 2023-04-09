# cook-import

A command line tool to import recipes into [Cooklang][3] format using [ChatGPT][4]

> :warning: Warning: This repository is in the middle of development and so the
> mode of operation may change at any time.

## TL;DR

Generate an [OpenAI API key][1].

```shell
git clone https://github.com/nicholaswilde/cook-import.git
cd cook-import
go run ./... -l <recipe link> -k <openai-api-key> -f
```

## Notes

`cook-import` uses the [viper][2] library and so environmental variables and a
config file may also be used.

[1]: <https://platform.openai.com/account/api-keys>
[2]: <https://github.com/spf13/viper>
[3]: <https://cooklang.org/>
[4]: <https://chat.openai.com/>
