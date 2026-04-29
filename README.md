
## 1. 快思考与慢思考——Two-Stage ReAct循环
“When tools are available, models tend to act quickly rather than think deeply.”（当工具可用时，模型倾向于迅速采取行动，而不是深入思考。）
- 好处：使得模型具备可观测性，强制模型进行思考，提供了一个介入点
- 问题：think trace带来的上下文堆积问题
## 2. 极简主义与YOLO (You Only Live Once)——大量工具带来的上下文膨胀问题
如果一个agent运行时附带了很多的工具，往往会带来上下文膨胀的问题，消耗token，还使得大模型的注意力很分散，维护这些工具也很困难。
- 解决：使用简单的几个工具，如write_file, read_file, bash工具，在workDir中给予agent最高的bash权限
