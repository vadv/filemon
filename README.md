# Описание

Утилита для тейлинга логов, которая фильтрует необходимые строки и собирает по ним полезную информацию, что бы быстро отдать.

# Запуск

* тестирует (отсылая ping и ожидая pong) доступность unix-сокета /tmp/<binname>-mon.sock

## Сокет не доступен

* при не доступности сокета пытается удалить его (при его наличии) и создать заново.
* демон в отдельном треде тестирует ping-pong и при недоступности - выходит.

## Сокет доступен

* клиент опрашивает сервер и получает ответ, если нет ошибок, в виде float, если есть ошибки они уходят в stderr, также индикатором является не нулевой выход клиента.

# Параметры запуска

`$ gofilemon <path-to-file> <regexp> <command> <command-args...>`

`path-to-file` - это путь до файла, при первом запросе файл должен быть доступен. отслеживается его truncate, delete

`regexp` - это фильтр для строк, все не подходящие строки игнорируются

`command` - вкомпиленный в бинарь сценарий сбора информации, на текущий момент:

* CountLine
* CountLinePerSecond
* SummField
* SummFieldPerSecond

`command-args` - дополнительные аргуемнты для сценария (например число -3 для SummField означает что суммировать надо третью с конца колонку, разделенные пробелами)
