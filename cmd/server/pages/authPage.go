package pages

// CreateAuthPage returns an HTML page that extracts access token from URL fragment
// and sends it to the server before redirecting to Discord
func CreateAuthPage() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>Discord Authentication</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin: 50px;
            background-color: #36393f;
            color: #ffffff;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #2f3136;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        }
        h1 {
            color: #7289da;
        }
        #status {
            margin: 20px 0;
            padding: 10px;
            border-radius: 4px;
        }
        .success {
            background-color: #43b581;
        }
        .error {
            background-color: #f04747;
        }
        .loading {
            background-color: #faa61a;
        }
        button {
            background-color: #7289da;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            margin-top: 20px;
        }
        button:hover {
            background-color: #677bc4;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Discord Authentication</h1>
        <div id="status" class="loading">Processing your authentication...</div>
        <div id="token-info"></div>
        <button id="redirect-btn" style="display:none;">Continue to Discord</button>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const statusEl = document.getElementById('status');
            const tokenInfoEl = document.getElementById('token-info');
            const redirectBtn = document.getElementById('redirect-btn');
            
            // Function to extract hash parameters
            function getHashParams() {
                const hash = window.location.hash.substring(1);
                const params = {};
                
                if (!hash) {
                    return params;
                }
                
                hash.split('&').forEach(pair => {
                    const [key, value] = pair.split('=');
                    params[key] = decodeURIComponent(value || '');
                });
                
                return params;
            }
            
            // Extract token from URL fragment
            const params = getHashParams();
            const accessToken = params['access_token'];
            const tokenType = params['token_type'];
            const expiresIn = params['expires_in'];
            
            if (!accessToken) {
                statusEl.textContent = 'Error: No access token found in URL';
                statusEl.className = 'error';
                return;
            }
            
            // Send token to server
            fetch('/store-token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    access_token: accessToken,
                    token_type: tokenType,
                    expires_in: expiresIn
                })
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to store token');
                }
                return response.json();
            })
            .then(data => {
                statusEl.textContent = 'Authentication successful!';
                statusEl.className = 'success';
                tokenInfoEl.innerHTML = '<p>Your token has been securely stored.</p>';
                
                // Show redirect button
                redirectBtn.style.display = 'inline-block';
                redirectBtn.addEventListener('click', function() {
                    window.location.href = 'https://discord.com/app';
                });
                
                // Auto redirect after 5 seconds
                setTimeout(() => {
                    window.location.href = 'https://discord.com/app';
                }, 5000);
            })
            .catch(error => {
                console.error('Error:', error);
                statusEl.textContent = 'Error: ' + error.message;
                statusEl.className = 'error';
            });
        });
    </script>
</body>
</html>
`
}
