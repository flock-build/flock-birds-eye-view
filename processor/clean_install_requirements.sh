# Uninstall all currently installed packages
pip freeze > all_packages.txt
pip uninstall -r all_packages.txt -y
rm all_packages.txt

# Install packages from requirements.txt
pip install -r requirements.txt