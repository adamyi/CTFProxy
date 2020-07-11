# NOTES(adamyi@):
# port python 3.8 html escape to python 2 cuz it's such a pain to set up python 3
# from https://github.com/python/cpython/blob/3.8/Lib/html/__init__.py
def escape(s, quote=True):
  """
    Replace special characters "&", "<" and ">" to HTML-safe sequences.
    If the optional flag quote is true (the default), the quotation mark
    characters, both double quote (") and single quote (') characters are also
    translated.
    """
  s = s.replace("&", "&amp;")  # Must be done first!
  s = s.replace("<", "&lt;")
  s = s.replace(">", "&gt;")
  if quote:
    s = s.replace('"', "&quot;")
    s = s.replace('\'', "&#x27;")
  return s
