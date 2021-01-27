import {Socket} from "./node_modules/phoenix/assets/js/phoenix";
import {LiveSocket, debug} from "./node_modules/phoenix_live_view/assets/js/phoenix_live_view";

let liveSocket = new LiveSocket("/live", Socket, {});

liveSocket.connect()

liveSocket.enableDebug()
