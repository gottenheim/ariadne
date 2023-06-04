# Project Ariadne
A spaced repetition approach to memorize programming languages and technologies.

## Quick overview
Please, check out [shell example](https://github.com/gottenheim/ariadne-shell-example) to have a quick hands-on experience with Ariadne tool.

## Reasons to create a new tool
Space repetition tools like SuperMemo and Anki are great to memorize facts if all you need to remember them is your memory and maybe a piece of paper. Things are getting harder if you need to remember new programming language or technology - the more experience you have, the more confident you are. I started my journey from creating Anki cards for different programming constructs but finally understood that I can't rely on my knowledge because of lacking experience. You can say - well, but you can answer "code" Anki questions by sitting opposite your code editor, typing code, compiling and running it. It's possible, but have some limitations. 

Imagine you are learning shell commands and you need to answer a question about special filtering with grep or awk. Perhaphs you will need a sample file to check your knowledge. Finding this file can take some time. You can then realize that it would be good to keep these files together and group them somehow. 

Now imagine you're learning something more sophisticated - for example, PostgreSQL commands, using [Postgres tutorial](https://www.postgresqltutorial.com/). To answer such cards, it would be better to have some temporary instance of PostgreSQL server, where you can quickly import test database to have possibility to restore it after harmful changes. Having such requirements, I usually start thinking about Docker and virtualization. Remember, we are still talking about answering one (maybe) simple question. 

So finally I realized that code cards and artifacts need to be stored together and managed with some scripting/coded environment, that will allow to create new cards from basic template, install dependencies (Docker, PostgreSQL, Minikube, etc), run code, save answer somewhere and show it on demand.




