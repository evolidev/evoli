
// import { createApp, reactive } from 'https://unpkg.com/petite-vue?module'
const { createApp, reactive } = PetiteVue;

window.onpopstate = function(event) {
    // alert("location: " + document.location + ", state: " + JSON.stringify(event.state));
    send({
        action: 'link',
        value: document.location.pathname,
    })
}

// document.addEventListener('click', function (e) {
//     console.log(parseOutMethodAndParams(
//         'Test(1,2,3, user.id)',
//         {user: {id: 9}, $event: e}
//     ));
// })

function parseOutMethodAndParams (rawMethod, scope) {
    let method = rawMethod
    let params = []
    const parts = method.match(/(.*?)\((.*)\)/s)

    if (parts) {
        method = parts[1]

        let func = new Function(Object.keys(scope), `return (function () {
            for (var l=arguments.length, p=new Array(l), k=0; k<l; k++) {
                p[k] = arguments[k];
            }
            return [].concat(p);
        })(${parts[2]})`)

        params = func(...Object.values(scope))
    }

    return {method, params}
}


const componentStates = reactive({});

window.Evo = {
    init:(id, data) => {
        // use this property to define component specific data
        const defaultComponent = {
            // name,
        }

        const dataObject = {
            ...defaultComponent,
            // ...BrightComponents[name] ?? {},
            ...componentStates[id] ?? {},
            ...data,
        };

        // console.info('Component INIT:', id, data, dataObject);

        return data.state;
    }
}

const onResponse = (data) => {
    // try to decode data
    try {
        data = JSON.parse(data);
    } catch (e) {
        console.error('Error parsing JSON', e);
    }

    if (!data.component?.name) {
        console.error('Invalid response', data);
        return;
    }

    if (!componentStates[data.component.name]) {
        componentStates[data.component.name] = {};
    }

    Object.assign(componentStates[data.component.name], data.state);

    if (data.type === 'page') {
        // replace the content of the page
        document.querySelector('.router-view').innerHTML = data.content;
        // re-initialize components
        init();
    }
}

// manipulate it here
// store.inc()

const getComponentData = (data) => {
    // get component data that don't start with $
    return Object.keys(data).filter(key => key[0] !== '$').reduce((acc, key) => {
        // if data is an object, get its data
        if (typeof data[key] === 'object') {
            acc[key] ={...data[key]}
        } else {
            acc[key] = data[key]
        }
        return acc
    }, {})

}

const click = (context) => {
    const { el, effect, exp, ctx } = context;
    const handler = () => {
        // console.log(ctx.scope, context);
        console.log('click: ' + exp);
        const state = getComponentData(ctx.scope);
        const {component} = state;
        delete state.component;

        send({
            action: 'click',
            value: exp,
            state,
            component
        });
    };

    effect(() => {
        // console.log('v-click: ' + exp);
        el.addEventListener('click', handler)
    })

    return () => {
        // console.log('remove click: ' + exp);
        el.removeEventListener('click', handler)
    }
}

const link = (context) => {
    const { el, effect, exp, ctx } = context;

    const handler = () => {
        // console.log(ctx.scope, context);
        console.log('====> link: ' + exp);
        const state = getComponentData(ctx.scope);
        const {component} = state;
        delete state.component;

        send({
            action: 'link',
            value: exp,
            state,
            component
        });

        history.pushState({}, null, exp)
    }

    effect(() => {
        // console.log('v-link: ' + exp);
        el.addEventListener('click', handler)
    })

    return () => {
        // console.log('remove link: ' + exp);
        el.removeEventListener('click', handler)
    }
}

const init = ()  => {
    // get all node elements with the attribute 'data-scope'
    const nodes = document.querySelectorAll('[data-scope]');
    console.log('Found component nodes:', nodes);

    // loop through all nodes and register the component
    for (let i = 0; i < nodes.length; i++) {
        const node = nodes[i];
        const id = node.getAttribute('data-cid');
        const dataRaw = node.getAttribute('data-scope');

        const data = JSON.parse(dataRaw);

        console.log('Registering component:', id, data);

        let state = reactive({});
        if (id) {
            state = reactive(Evo.init(id, data));
        }

        componentStates[id] = state;

        console.log('Component state:', state, node);

        createApp(state)
            .directive('click', click)
            .directive('link', link)
            .mount(node);

        node.removeAttribute('data-scope');
    }
}

init();

window.init = init;

// createApp({store})
//     .directive('click', click)
//     .mount()

/**
 * Socket section
 * @type {HTMLElement}
 */
// let input = document.getElementById("input");
// let output = document.getElementById("output");
// let host = document.location.host;
// let socket = new WebSocket(`ws://${host}/component/pipeline`);
//
// socket.onopen = function () {
//     output.innerHTML += "Status: Connected\n";
// };
//
// socket.onmessage = function (e) {
//     output.innerHTML += "Server: " + e.data + "\n";
//     onResponse(e.data);
// };

// function submit() {
//     send({ input: input.value });
// }

function send(data) {
    console.log('SEND socket data', data);
    socket.send(JSON.stringify(data));
    // input.value = "";
}
