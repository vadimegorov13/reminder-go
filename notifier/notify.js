const notifier = require("node-notifier")
const path = require("path")

module.exports = {
  notify: ({ title, message }, callBack) => {
    notifier.notify(
        {
            title: title || "Empty title",
            message: message || "Empty message",
            sound: "Frog",
            wait: true,
            reply: true,
            closedLabel: "Completed?",
            timeout: 15,
            icon: path.join(__dirname, 'Ame.jpg'),
            
        },
        (err, response, reply) => {
            callBack(reply)
        }
    );

    notifier.on('click', () => {
      opn('', { app: 'discord' });
    });
    
    notifier.on('timeout', function (notifierObject, options) {
      console.log("timeout")
    });
  },
};
