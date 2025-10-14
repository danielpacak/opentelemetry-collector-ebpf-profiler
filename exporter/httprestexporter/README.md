# HTTP REST API Exporter

``` console
$ curl http://localhost:7799/workload/826396f...b5/functions | jq -r  '.[] | .File + " " + .Name' | sort
 do_futex
 do_syscall_64
 entry_SYSCALL_64_after_hwframe
 futex_wake
 _raw_spin_unlock_irqrestore
<string> _generated_cache_key_traversal
 try_to_wake_up
/usr/local/lib/python3.10/collections/__init__.py __getitem__
/usr/local/lib/python3.10/http/server.py handle
/usr/local/lib/python3.10/http/server.py handle_one_request
/usr/local/lib/python3.10/site-packages/engineio/middleware.py __call__
/usr/local/lib/python3.10/site-packages/flask/app.py __call__
/usr/local/lib/python3.10/site-packages/flask/app.py dispatch_request
/usr/local/lib/python3.10/site-packages/flask/app.py full_dispatch_request
/usr/local/lib/python3.10/site-packages/flask/app.py update_template_context
/usr/local/lib/python3.10/site-packages/flask/app.py wsgi_app
/usr/local/lib/python3.10/site-packages/flask_login/utils.py decorated_view
/usr/local/lib/python3.10/site-packages/flask_login/utils.py _get_user
/usr/local/lib/python3.10/site-packages/flask_login/utils.py _user_context_processor
/usr/local/lib/python3.10/site-packages/flask_security/core.py verify_and_update_password
/usr/local/lib/python3.10/site-packages/flask_security/decorators.py decorated
/usr/local/lib/python3.10/site-packages/flask_security/forms.py validate
/usr/local/lib/python3.10/site-packages/flask_security/utils.py default_render_template
/usr/local/lib/python3.10/site-packages/flask_security/utils.py verify_and_update_password
/usr/local/lib/python3.10/site-packages/flask_security/views.py login
/usr/local/lib/python3.10/site-packages/flask_socketio/__init__.py __call__
/usr/local/lib/python3.10/site-packages/flask/templating.py _render
/usr/local/lib/python3.10/site-packages/flask/templating.py render_template
```

``` console
$ curl http://localhost:7799/workload/b9a01e3...ff/process/2981/functions | \
  jq -r '.[] | select(.Language == "go" or .Language == "kernel") | .Language + ": " + .Name' | sort
go: bufio.(*Writer).Flush
go: crypto/ecdh.(*PrivateKey).PublicKey.func1
go: crypto/ecdh.(*x25519Curve).ecdh
go: crypto/ecdh.(*x25519Curve).privateKeyToPublicKey
go: crypto/ecdh.x25519ScalarMult
go: crypto.Hash.New
go: crypto/internal/bigmod.addMulVVW1024
go: crypto/internal/bigmod.(*Nat).shiftIn
go: crypto/internal/edwards25519/field.(*Element).carryPropagateGeneric
go: crypto/internal/edwards25519/field.(*Element).Swap
go: crypto/internal/mlkem768.nttMul
go: crypto/internal/mlkem768.pkeEncrypt
go: crypto/rsa.decrypt
go: crypto/rsa.(*PrivateKey).Sign
go: crypto/rsa.SignPSS
go: crypto/rsa.signPSSWithSalt
go: crypto/rsa.VerifyPKCS1v15
go: crypto/sha256.New
go: crypto/tls.(*Conn).serverHandshake-fm
go: crypto/tls.(*Conn).write
go: crypto/tls.(*Conn).Write
go: crypto/x509.(*CertPool).findPotentialParents
go: encoding/json.(*Decoder).Decode
go: encoding/json.(*Decoder).readValue
go: encoding/json.Unmarshal
go: github.com/emicklei/go-restful/v3.(*Container).dispatch
go: github.com/felixge/httpsnoop.(*rw).WriteHeader
go: github.com/prometheus/client_golang/prometheus.(*histogram).observe
go: github.com/prometheus/client_golang/prometheus.(*histogram).Observe
```
