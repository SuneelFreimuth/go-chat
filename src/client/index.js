function main() {
    const chat = new Chat(document.querySelector('.chat-container'))
    const socket = new Socket('ws://localhost:5000/ws', msg => {
        chat.appendMessage(msg)
    })
    chat.onSubmit(body => {
        // TODO: Author name
        socket.send(new MessageOut('Someone', body))
    })
}

class Socket {
    constructor(url, onMessage) {
        this.url = url
        this.onMessage = onMessage
        this.ws = new WebSocket(this.url)
        this.ws.onopen = e => {
            console.log('Established connection with server!')
        }
        this.ws.onmessage = this.#onMessage
    }

    #onMessage = e => {
        const {author, body, serverTime} = JSON.parse(e.data)
        console.log('Received message:', author, body, serverTime)
        this.onMessage(new MessageIn(author, body, serverTime))
    }

    send = body => {
        console.log('Sending message:', body)
        this.ws.send(JSON.stringify(body))
    }
}

class Chat {
    /**
     * @param {HTMLDivElement} container 
     */
    constructor(container) {
        this.root = container
        this.messages = container.children[0]
        this.input = container.children[1]
        this._onSubmit = () => {}
    }

    /** @param {MessageIn} msg */
    appendMessage = msg => {
        const author = document.createElement('p')
        author.classList.add('author')
        author.innerText = msg.author

        const body = document.createElement('p')
        body.classList.add('body')
        body.innerText = msg.body

        const message = document.createElement('div')
        message.classList.add('message')
        message.appendChild(author)
        message.appendChild(body)

        this.messages.appendChild(message)
    }

    onSubmit = handler => {
        this.input.removeEventListener('keyup', this._onSubmit)
        this._onSubmit = e => {
            if (e.key === 'Enter') {
                handler(e.target.value)
                e.target.value = ''
            }
        }
        this.input.addEventListener('keyup', this._onSubmit)
    }
}

class MessageIn {
    constructor(author, body) {
        // TODO: Include author.
        this.author = "Someone"
        this.body = body
    }
}

class MessageOut {
    constructor(author, body, serverTime) {
        this.author = author
        this.body = body
        this.serverTime = serverTime
    }
}

main()