# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'CCNP'
copyright = '2023, Intel'
author = 'Ken Lu, Ruoyu Ying, Hairong Chen, Xiaocheng Dong, Lei Zhou'

import os
import sys
sys.path.insert(0, os.path.abspath('../sdk/python3/ccnp/eventlog'))
sys.path.insert(0, os.path.abspath('../sdk/python3/ccnp/measurement'))
sys.path.insert(0, os.path.abspath('../sdk/python3/ccnp/quote'))

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = ['sphinx.ext.autodoc', 'sphinx.ext.napoleon', 'sphinx_markdown_builder',
        'sphinx_mdinclude']

napoleon_google_docstring = True

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

import sphinx_rtd_theme

html_theme = "sphinx_rtd_theme"
html_theme_path = [sphinx_rtd_theme.get_html_theme_path()]
