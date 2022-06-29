
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

function parseOutMethodAndParams (rawMethod, scope = {}) {
    let method = rawMethod
    let parameters = []
    const parts = method.match(/(.*?)\((.*)\)/s)

    if (parts) {
        method = parts[1]

        let func = new Function(Object.keys(scope), `return (function () {
            for (var l=arguments.length, p=new Array(l), k=0; k<l; k++) {
                p[k] = arguments[k];
            }
            return [].concat(p);
        })(${parts[2]})`)

        parameters = func(...Object.values(scope))
    }

    return {method, parameters}
}


const componentStates = reactive({});
const components = {};

window.Evo = {
    init:(data) => {
        // use this property to define component specific data
        const defaultComponent = {
            // name,
        }

        let {state, _id, component} = data
        state = reactive({...state, _id})

        // if (componentStates[_id]) {
        //     // state = {...componentStates[_id], ...state}
        //     Object.assign(componentStates[_id], state)
        // } else {

        // }

        // loop through state and make everything reactive
        // for (let key in state) {
        //     state[key] = reactive(state[key])
        // }

        components[_id] = { ...data }

        componentStates[_id] = state

        console.log('Init component:', data, componentStates[_id]);

        return componentStates[_id]
    }
}



const onResponse = (data) => {
    // try to decode data
    let response
    try {
        response = JSON.parse(data)
    } catch (e) {
        console.error('Error parsing JSON', e)
        return
    }

    console.log('Received response:', response);


    if (!response.component) {
        console.error('Invalid response', data)
        return;
    }

    const { state, _id, component } = response


    if (!componentStates[_id]) {
        componentStates[_id] = {}
    }

    console.log(componentStates[_id], state);

    Object.assign(componentStates[_id], state)

    // for (let key in state) {
    //     componentStates[_id][key] = state[key]
    // }

    if (response.type === 'page') {
        // replace the content of the page
        document.querySelector('.router-view').innerHTML = response.content;
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
        const state = getComponentData(ctx.scope);
        // console.log(ctx.scope, context);
        // console.log('click: ' + exp, state, context);
        // const {component} = state;
        // delete state.component;

        send({
            action: 'click',
            expression: exp,
            state,
            // component
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

// const init = ()  => {
//     // get all node elements with the attribute 'v-scope'
//     const nodes = document.querySelectorAll('[v-scope]');
//     console.log('Found component nodes:', nodes);
//
//     // loop through all nodes and register the component
//     for (let i = 0; i < nodes.length; i++) {
//         const node = nodes[i];
//         const id = node.getAttribute('data-cid');
//         const dataRaw = node.getAttribute('v-scope');
//
//         const data = JSON.parse(dataRaw);
//
//         console.log('Registering component:', id, data);
//
//         let state = reactive({});
//         if (id) {
//             state = reactive(Evo.init(id, data));
//         }
//
//         componentStates[id] = state;
//
//         console.log('Component state:', state, node);
//
//         createApp(state)
//             .directive('click', click)
//             .directive('link', link)
//             .mount(node);
//
//         node.removeAttribute('v-scope');
//     }
// }

const play = reactive({
    Name: 'OMer',
    Email: 'OMer',
    Type: 'OMer',
    Password: 'OMer',
    mount(data) {
        console.log('mounting', play);
        return play
        return Evo.init(data)
    }
})

// init();
const app = PetiteVue
    .createApp({
        mount(data) {
            // return play
            return Evo.init(data)
        }
    })
    // .createApp(play)
    .directive('click', click)
    .directive('link', link)
    .mount();

// window.init = init;

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
    console.log('send', data);

    const requestData = parseOutMethodAndParams(data.expression, data.state)
    console.log(requestData, data);

    let state = {...data.state}
    const id = state._id

    // remove _id from state
    delete state._id

    const component = components[id]
    console.log('Component:', component);

    const request = {
        _id: component._id,
        component: component.component,
        state,
        action: data.action,
        ...requestData
    }

    sendXhr(request)

    // sendXhr({
    //     component: 'Login',
    //     state: {},
    //     action: 'click',
    // })

    // console.log('SEND socket data', data);
    // socket.send(JSON.stringify(data));
    // input.value = "";
}

function sendXhr(data) {
    console.log('SEND xhr data', data);
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/internal/component');
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(data));

    // handle response
    xhr.onload = function() {
        onResponse(xhr.responseText);
    }
}