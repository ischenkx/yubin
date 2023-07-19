import axios, {Axios} from "axios";
import PubSub from "../controllers/pubsub";
import Sources from "../controllers/sources";
import Templates from "../controllers/templates";
import Users from "../controllers/users";

export default class Client {
    private client: Axios
    private _pubsub: PubSub
    private _sources: Sources
    private _templates: Templates
    private _users: Users

    constructor(url: string) {
        this.client = axios.create({
            baseURL: url,
        })

        this._pubsub = new PubSub(this.client)
        this._sources = new Sources(this.client)
        this._templates = new Templates(this.client)
        this._users = new Users(this.client)
    }

    public pubsub(): PubSub {
        return this._pubsub
    }

    public sources(): Sources {
        return this._sources
    }

    public templates(): Templates {
        return this._templates
    }

    public users(): Users {
        return this._users
    }
}