let app = {
    init: function () {
        console.log("Wait for the ready signal ...")
        // Wait for the ready signal
        document.addEventListener('astilectron-ready', function () {
            // This will send a message to GO
            astilectron.sendMessage("hello", function (message) {
                console.log("received " + message)
            });
        });
    }
}