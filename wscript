#

top = '.'
out = '__build__'

def options(ctx):
    ctx.load('go')
    pass

def configure(ctx):
    ctx.load('go')
    ctx.load('gcc')
    ctx.load('c')
    
def build(ctx):

    ctx(
        features='go gopackage',
        name='go-hep/io',
        source='pkg/hep/io/io.go',
        target='hep/io',
        )

    ctx(
        features='go gopackage',
        name='go-hep/fourmom',
        source='pkg/hep/fourmom/fourmom.go',
        target='hep/fourmom',
        )

    ctx(
        features='go gopackage',
        name='go-hep/pdg',
        source='pkg/hep/pdg/pdg.go',
        target='hep/pdg',
        )
        
        
