#!/bin/bash


generateKeys(){
    # Generate private key
    openssl genpkey -algorithm RSA -out private_key.pem

    # Extract public key from private key and base64 encode
    openssl rsa -pubout -in private_key.pem | base64 -w 0 > public_key.base64

    # Base64 encode private key
    base64 -w 0 private_key.pem > private_key.base64

    # Create .env file and store encoded keys
    echo "JWT_PRIVATE_KEY=$(cat public_key.base64)" > .env
    echo "" >> .env
    echo "JWT_PUBLIC_KEY=$(cat private_key.base64)" >> .env

    # Clean up temporary files
    rm public_key.base64 private_key.base64 private_key.pem

    echo "generated keys"
}

generateKeys


replace_module_name() {
    echo "Enter your module name (e.g., github.com/yourusername/yourmodule):"
    read new_module_name

    if [[ -z "$new_module_name" ]]; then
        echo "Module name cannot be empty."
        exit 1
    fi

    # Replace old module name with the new one in all files
    find . -type f -exec sed -i "s/github\.com\/tomdoestech\/goth/$new_module_name/g" {} +

    echo "Module name replaced in all files."
}

replace_module_name